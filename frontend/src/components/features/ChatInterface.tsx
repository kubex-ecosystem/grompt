import { Loader2, Send, Sparkles, User } from 'lucide-react';
import React, { FormEvent, useMemo, useState } from 'react';
import Card from '../ui/Card';

type ChatRole = 'user' | 'assistant' | 'system';

interface ChatMessage {
  id: string;
  role: ChatRole;
  content: string;
  createdAt: number;
  usedProvider?: string;
}

interface ChatInterfaceProps {
  onSend?: (messages: ChatMessage[], input: string) => Promise<{ content: string; provider?: string } | null>;
}

const ChatBubble: React.FC<{ message: ChatMessage }> = ({ message }) => {
  const isUser = message.role === 'user';
  return (
    <div
      className={`max-w-2xl rounded-2xl border px-4 py-3 text-sm shadow-sm transition ${
        isUser
          ? 'ml-auto border-[#06b6d4] bg-[#06b6d4] text-white shadow-soft-card'
          : 'mr-auto border-slate-200/80 bg-white text-[#475569] dark:border-[#13263a]/80 dark:bg-[#0a1523]/80 dark:text-[#e5f2f2]'
      }`}
    >
      <div className="flex items-center justify-between gap-3">
        <span className="inline-flex items-center gap-2 text-xs uppercase tracking-wide text-[#94a3b8] dark:text-[#64748b]">
          {isUser ? <User className="h-3.5 w-3.5" /> : <Sparkles className="h-3.5 w-3.5" />}
          {isUser ? 'Você' : 'Assistant'}
        </span>
        {message.usedProvider && (
          <span className="text-[10px] uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
            {message.usedProvider}
          </span>
        )}
      </div>
      <p className="mt-2 whitespace-pre-wrap leading-relaxed">{message.content}</p>
      <span className="mt-3 block text-[10px] uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
        {new Date(message.createdAt).toLocaleTimeString()}
      </span>
    </div>
  );
};

const ChatInterface: React.FC<ChatInterfaceProps> = ({ onSend }) => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [isSending, setIsSending] = useState(false);
  const disabled = input.trim().length === 0 || isSending;

  const placeholder = useMemo(
    () =>
      'Descreva o contexto do atendimento, insira scripts de vendas ou cole uma conversa existente para gerar respostas contextualizadas.',
    []
  );

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const trimmed = input.trim();
    if (!trimmed) return;

    const nextMessage: ChatMessage = {
      id: `msg-${Date.now()}`,
      role: 'user',
      content: trimmed,
      createdAt: Date.now(),
    };
    setMessages((prev) => [...prev, nextMessage]);
    setInput('');

    if (!onSend) return;

    setIsSending(true);
    try {
      const response = await onSend([...messages, nextMessage], trimmed);
      if (response) {
        const assistantMessage: ChatMessage = {
          id: `assistant-${Date.now()}`,
          role: 'assistant',
          content: response.content,
          createdAt: Date.now(),
          usedProvider: response.provider,
        };
        setMessages((prev) => [...prev, assistantMessage]);
      }
    } catch (error) {
      const assistantMessage: ChatMessage = {
        id: `assistant-${Date.now()}`,
        role: 'assistant',
        content:
          error instanceof Error
            ? `Não foi possível obter uma resposta: ${error.message}`
            : 'Não foi possível obter uma resposta da IA.',
        createdAt: Date.now(),
      };
      setMessages((prev) => [...prev, assistantMessage]);
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className="space-y-6">
      <Card title="Chat assistido" description="Converse com seu provedor principal usando contexto governado">
        <div className="flex flex-col gap-4">
          <div className="space-y-3">
            {messages.length === 0 ? (
              <p className="rounded-2xl border border-dashed border-slate-200/80 bg-white/70 p-6 text-sm text-[#64748b] dark:border-[#13263a]/70 dark:bg-[#0a1523]/60 dark:text-[#94a3b8]">
                Nenhuma mensagem ainda. Use o campo abaixo para iniciar uma conversa.
              </p>
            ) : (
              <div className="space-y-4">
                {messages.map((message) => (
                  <ChatBubble key={message.id} message={message} />
                ))}
              </div>
            )}
          </div>
          <form onSubmit={handleSubmit} className="rounded-2xl border border-slate-200/80 bg-white/85 p-4 shadow-sm dark:border-[#13263a]/70 dark:bg-[#0a1523]/70">
            <label htmlFor="chat-input" className="mb-2 block text-xs font-semibold uppercase tracking-[0.4em] text-[#94a3b8] dark:text-[#64748b]">
              Sua mensagem
            </label>
            <textarea
              id="chat-input"
              value={input}
              onChange={(event) => setInput(event.target.value)}
              placeholder={placeholder}
              rows={4}
              className="w-full resize-none rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2] dark:focus:border-[#38cde4] dark:focus:ring-[#38cde4]/20"
            />
            <div className="mt-3 flex items-center justify-between">
              <p className="text-[11px] uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                Histórico salvo localmente
              </p>
              <button
                type="submit"
                disabled={disabled}
                className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-5 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#0891b2] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/30 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
              >
                {isSending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
                {isSending ? 'Gerando...' : 'Enviar mensagem'}
              </button>
            </div>
          </form>
        </div>
      </Card>
    </div>
  );
};

export type { ChatMessage };
export default ChatInterface;
