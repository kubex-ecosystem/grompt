package types

import (
	"bufio"
	"encoding/asn1"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	ci "github.com/kubex-ecosystem/grompt/internal/interfaces"
	gl "github.com/kubex-ecosystem/logz/logger"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Mapper serializa/desserializa T com decoders streaming.
// Contratos:
// - JSON: se T for slice, aceita array JSON ou sequência de objetos; se T não for slice, 1 objeto.
// - YAML: se T for slice, aceita múltiplos docs (---); se T não for slice, 1 doc.
// - TOML: 1 doc.
// - XML: 1 doc.
// - ENV: apenas map[string]string.
type Mapper[T any] struct {
	FilePath string
	Ptr      *T
}

func NewMapper[T any](object *T, filePath string) ci.IMapper[T] {
	return &Mapper[T]{FilePath: filePath, Ptr: object}
}

func NewMapperW[T any](object T, filePath string) ci.IMapper[T] {
	return &Mapper[T]{FilePath: filePath, Ptr: &object}
}

func NewMapperFromFile[T any](filePath string) ci.IMapper[T] {
	var obj T
	return &Mapper[T]{FilePath: filePath, Ptr: &obj}
}

func NewMapperFromFileW[T any](filePath string) ci.IMapper[T] {
	var obj T
	return &Mapper[T]{FilePath: filePath, Ptr: &obj}
}

func NewMapperFromObject[T any](object *T) ci.IMapper[T] {
	return &Mapper[T]{Ptr: object}
}

func NewMapperFromObjectW[T any](object T) ci.IMapper[T] {
	return &Mapper[T]{Ptr: &object}
}

// -------------------- Helpers --------------------

func detectFormatByExt(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".xml":
		return "xml"
	case ".toml", ".tml":
		return "toml"
	case ".env":
		return "env"
	default:
		return ""
	}
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// -------------------- Serialize --------------------

func (m *Mapper[T]) Serialize(format string) ([]byte, error) {
	if m.Ptr == nil {
		return nil, errors.New("mapper: ponteiro de destino nil")
	}
	switch strings.ToLower(format) {
	case "json", "js":
		return json.Marshal(m.Ptr)
	case "yaml", "yml":
		return yaml.Marshal(m.Ptr)
	case "xml", "html":
		return xml.Marshal(m.Ptr)
	case "toml", "tml":
		return toml.Marshal(m.Ptr)
	case "asn", "asn1":
		return asn1.Marshal(*m.Ptr)
	case "env", "envs", ".env", "environment":
		// Apenas map[string]string é suportado para ENV
		rv := reflect.ValueOf(m.Ptr).Elem()
		if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String && rv.Type().Elem().Kind() == reflect.String {
			env := make(map[string]string, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				env[iter.Key().String()] = iter.Value().String()
			}
			// gotenv não tem Encoder streaming; monta manualmente.
			var b strings.Builder
			for k, v := range env {
				// simples; se quiser suportar escaping avançado, tratar aqui.
				fmt.Fprintf(&b, "%s=%s\n", k, v)
			}
			return []byte(b.String()), nil
		}
		return nil, fmt.Errorf("ENV exige map[string]string; recebido: %T", *m.Ptr)
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

func (m *Mapper[T]) SerializeToFile(format string) {
	if format == "" {
		format = detectFormatByExt(m.FilePath)
	}
	data, err := m.Serialize(format)
	if err != nil {
		gl.Log("error", fmt.Sprintf("Error serializing object: %v", err))
		return
	}
	if err := ensureParentDir(m.FilePath); err != nil {
		gl.Log("error", fmt.Sprintf("Error creating parent dir: %v", err))
		return
	}
	f, err := os.OpenFile(m.FilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		gl.Log("error", fmt.Sprintf("Error opening file: %v", err))
		return
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			gl.Log("error", fmt.Sprintf("Error closing file: %v", cerr))
		}
	}()
	if _, err := f.Write(data); err != nil {
		gl.Log("error", fmt.Sprintf("Error writing file: %v", err))
		return
	}
	gl.Log("debug", fmt.Sprintf("Serialized to %s (%s) [%d bytes]", m.FilePath, strings.ToUpper(format), len(data)))
}

// -------------------- Deserialize (streaming) --------------------

func (m *Mapper[T]) Deserialize(object []byte, format string) (*T, error) {
	if m.Ptr == nil {
		return nil, errors.New("mapper: ponteiro de destino nil")
	}
	var r io.Reader = strings.NewReader(string(object))
	switch strings.ToLower(format) {
	case "json", "js":
		return m.decodeJSONStream(r)
	case "yaml", "yml":
		return m.DecodeYAMLStream(r)
	case "xml", "html":
		return m.DecodeXMLStream(r)
	case "toml", "tml":
		return m.DecodeTOMLStream(r)
	case "env", "envs", ".env", "environment":
		return m.DecodeENVStream(r)
	case "asn", "asn1":
		// ASN.1 não tem decoder streaming idiomático em std; teria que ler bytes.
		// Como geralmente é pequeno para config, leremos via io.ReadAll? Evitamos: então rejeita aqui.
		return nil, errors.New("ASN.1 streaming não suportado no momento")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

func (m *Mapper[T]) DeserializeFromFile(format string) (*T, error) {
	if m.Ptr == nil {
		return nil, errors.New("mapper: ponteiro de destino nil")
	}
	if _, err := os.Stat(m.FilePath); err != nil {
		gl.Log("error", fmt.Sprintf("File does not exist: %v", err))
		return nil, err
	}
	f, err := os.Open(m.FilePath)
	if err != nil {
		gl.Log("error", fmt.Sprintf("Error opening file: %v", err))
		return nil, err
	}
	defer func() {
		gl.Log("debug", "Closing input file")
		if cerr := f.Close(); cerr != nil {
			gl.Log("error", fmt.Sprintf("Error closing file: %v", cerr))
		}
	}()

	if format == "" {
		format = detectFormatByExt(m.FilePath)
		if format == "" {
			return nil, fmt.Errorf("não foi possível detectar o formato pelo sufixo de %s; informe o formato", m.FilePath)
		}
	}
	switch strings.ToLower(format) {
	case "json", "js":
		return m.decodeJSONStream(f)
	case "yaml", "yml":
		return m.DecodeYAMLStream(f)
	case "xml", "html":
		return m.DecodeXMLStream(f)
	case "toml", "tml":
		return m.DecodeTOMLStream(f)
	case "env", "envs", ".env", "environment":
		return m.DecodeENVStream(f)
	case "asn", "asn1":
		// ASN.1 não tem decoder streaming idiomático em std; teria que ler bytes.
		// Como geralmente é pequeno para config, leremos via io.ReadAll? Evitamos: então rejeita aqui.
		return nil, errors.New("ASN.1 streaming não suportado no momento")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

func (m *Mapper[T]) decodeJSONStream(r io.Reader) (*T, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields() // opcional; ajuda a pegar chaves erradas
	rt := reflect.TypeOf(*m.Ptr)
	if rt.Kind() == reflect.Slice {
		// T = []E
		elemType := rt.Elem()
		sliceVal := reflect.MakeSlice(rt, 0, 0)

		// Suporta array JSON OU sequência de objetos
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("JSON: erro lendo primeiro token: %w", err)
		}
		if delim, ok := tok.(json.Delim); ok && delim == '[' {
			// [ obj, obj, ... ]
			for dec.More() {
				elemPtr := reflect.New(elemType).Interface()
				if err := dec.Decode(elemPtr); err != nil {
					return nil, fmt.Errorf("JSON: erro decodificando elemento do array: %w", err)
				}
				sliceVal = reflect.Append(sliceVal, reflect.ValueOf(elemPtr).Elem())
			}
			// consumir ']'
			if _, err := dec.Token(); err != nil {
				return nil, fmt.Errorf("JSON: erro lendo token de fechamento do array: %w", err)
			}
		} else {
			// não começou com '[', então assumimos sequência de objetos
			dec.Buffered() // noop; apenas semântica
			// o token lido (tok) pode ser '{' já consumido; então usar Decode direto no primeiro elemento
			// Estratégia: retroceder não dá; então decodificar o "restante" como primeiro elem + laço More() falso.
			// Para simplificar: usamos um decoder sobre um reader que já consumiu um token.
			// Na prática, o token foi '{', então o próximo Decode pega o objeto inteiro.
			firstElem := reflect.New(elemType).Interface()
			if err := dec.Decode(firstElem); err != nil {
				return nil, fmt.Errorf("JSON: erro decodificando primeiro elemento fora de array: %w", err)
			}
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(firstElem).Elem())

			// Tentar mais objetos sequenciais (válido para NDJSON/concatenados)
			for {
				// pular espaços/brancos
				if err := consumeWhitespace(dec); err != nil && !errors.Is(err, io.EOF) {
					return nil, fmt.Errorf("JSON: erro consumindo whitespace: %w", err)
				}
				// Tentar próximo objeto
				nextElem := reflect.New(elemType).Interface()
				if err := dec.Decode(nextElem); err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					// Pode ter acabado exatamente; se for “invalid character ‘}’ looking for beginning of value”
					// porque não há próximo objeto, encerramos.
					if isBenignJSONEnd(err) {
						break
					}
					return nil, fmt.Errorf("JSON: erro decodificando elemento subsequente: %w", err)
				}
				sliceVal = reflect.Append(sliceVal, reflect.ValueOf(nextElem).Elem())
			}
		}
		reflect.ValueOf(m.Ptr).Elem().Set(sliceVal)
		return m.Ptr, nil
	}

	// T = objeto/map simples
	if err := dec.Decode(m.Ptr); err != nil {
		return nil, fmt.Errorf("JSON: erro decodificando objeto: %w", err)
	}
	// garantir que não há lixo após o objeto
	if err := ensureJSONEOF(dec); err != nil {
		return nil, err
	}
	return m.Ptr, nil
}

func consumeWhitespace(dec *json.Decoder) error {
	// json.Decoder não expõe “skip”; a próxima Decode já ignora whitespace.
	// Então não há nada pra fazer de fato; retornamos nil.
	return nil
}

func isBenignJSONEnd(err error) bool {
	// Heurística: erros comuns ao terminar sequência sem próximo valor
	msg := err.Error()
	return strings.Contains(msg, "looking for beginning of value") ||
		strings.Contains(msg, "unexpected EOF") ||
		strings.Contains(msg, "EOF")
}

func ensureJSONEOF(dec *json.Decoder) error {
	// tentar consumir espaços e verificar se acabou
	var tmp any
	if err := dec.Decode(&tmp); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		// Se deu erro porque não há mais um objeto válido, isso é OK
		if isBenignJSONEnd(err) {
			return nil
		}
		return fmt.Errorf("JSON: dados extras após o objeto: %w", err)
	}
	return fmt.Errorf("JSON: dados extras após o objeto")
}

func (m *Mapper[T]) DecodeYAMLStream(r io.Reader) (*T, error) {
	dec := yaml.NewDecoder(r)
	rt := reflect.TypeOf(*m.Ptr)
	if rt.Kind() == reflect.Slice {
		elemType := rt.Elem()
		sliceVal := reflect.MakeSlice(rt, 0, 0)
		for {
			elemPtr := reflect.New(elemType).Interface()
			if err := dec.Decode(elemPtr); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, fmt.Errorf("YAML: erro decodificando documento: %w", err)
			}
			// documentos vazios podem vir como zero value; ainda assim anexamos
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(elemPtr).Elem())
		}
		reflect.ValueOf(m.Ptr).Elem().Set(sliceVal)
		return m.Ptr, nil
	}

	// um único doc
	if err := dec.Decode(m.Ptr); err != nil {
		return nil, fmt.Errorf("YAML: erro decodificando: %w", err)
	}
	return m.Ptr, nil
}

func (m *Mapper[T]) DecodeTOMLStream(r io.Reader) (*T, error) {
	dec := toml.NewDecoder(r)
	if err := dec.Decode(m.Ptr); err != nil {
		return nil, fmt.Errorf("TOML: erro decodificando: %w", err)
	}
	return m.Ptr, nil
}

func (m *Mapper[T]) DecodeXMLStream(r io.Reader) (*T, error) {
	dec := xml.NewDecoder(r)
	if err := dec.Decode(m.Ptr); err != nil {
		return nil, fmt.Errorf("XML: erro decodificando: %w", err)
	}
	return m.Ptr, nil
}

func (m *Mapper[T]) DecodeENVStream(r io.Reader) (*T, error) {
	// Apenas map[string]string
	rv := reflect.ValueOf(m.Ptr).Elem()
	if !(rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String && rv.Type().Elem().Kind() == reflect.String) {
		return nil, fmt.Errorf("ENV exige destino map[string]string; recebido: %T", *m.Ptr)
	}
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	sc := bufio.NewScanner(r)
	// Aumenta buffer para valores grandes (default é 64KB)
	const maxToken = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, 64*1024)
	sc.Buffer(buf, maxToken)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		// Suporta KEY=VALUE (pegar apenas primeiro '=')
		k, v, ok := cutOnce(line, '=')
		if !ok {
			// Linha inválida no contexto .env — falha explícita (ou poderíamos ignorar)
			return nil, fmt.Errorf("ENV: linha inválida (esperado KEY=VALUE): %q", line)
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		// Remover aspas externas simples/duplas
		v = strings.Trim(v, `"'`)
		rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("ENV: erro lendo arquivo: %w", err)
	}
	return m.Ptr, nil
}

func cutOnce(s string, sep byte) (head, tail string, ok bool) {
	if i := strings.IndexByte(s, sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return "", "", false
}

// -------------------- Helpers mantidas --------------------

func SanitizeQuotesAndSpaces(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "'", "\"")
	input = strings.Trim(input, "\"")
	return input
}

func IsEqual(a, b string) bool {
	a, b = SanitizeQuotesAndSpaces(a), SanitizeQuotesAndSpaces(b)
	ptsEqual := levenshtein(a, b)
	maxLen := maxL(len(a), len(b))
	threshold := maxLen / 4
	return ptsEqual <= threshold
}

func maxL(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func levenshtein(s, t string) int {
	m, n := len(s), len(t)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}
	prevRow := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prevRow[j] = j
	}
	for i := 1; i <= m; i++ {
		currRow := make([]int, n+1)
		currRow[0] = i
		for j := 1; j <= n; j++ {
			cost := 1
			if s[i-1] == t[j-1] {
				cost = 0
			}
			a := prevRow[j] + 1
			b := currRow[j-1] + 1
			c := prevRow[j-1] + cost
			currRow[j] = min(a, b, c)
		}
		prevRow = currRow
	}
	return prevRow[n]
}
