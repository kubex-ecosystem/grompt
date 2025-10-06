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
          ? 'ml-auto bg-slate-900 text-white border-slate-900/80 shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)]'
          : 'mr-auto bg-white/90 border-slate-200/80 text-slate-700 dark:bg-[#10151b]/80 dark:border-slate-800/60 dark:text-slate-200'
      }`}
    >
      <div className="flex items-center justify-between gap-3">
        <span className="inline-flex items-center gap-2 text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">
          {isUser ? <User className="h-3.5 w-3.5" /> : <Sparkles className="h-3.5 w-3.5" />}
          {isUser ? 'Você' : 'Assistant'}
        </span>
        {message.usedProvider && (
          <span className="text-[10px] uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
            {message.usedProvider}
          </span>
        )}
      </div>
      <p className="mt-2 whitespace-pre-wrap leading-relaxed">{message.content}</p>
      <span className="mt-3 block text-[10px] uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
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
              <p className="rounded-2xl border border-dashed border-slate-200/80 bg-white/60 p-6 text-sm text-slate-500 dark:border-slate-700/80 dark:bg-[#0a0f14]/60 dark:text-slate-400">
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
          <form onSubmit={handleSubmit} className="rounded-2xl border border-slate-200/80 bg-white/80 p-4 shadow-sm dark:border-slate-800/60 dark:bg-[#0a0f14]/70">
            <label htmlFor="chat-input" className="mb-2 block text-xs font-semibold uppercase tracking-[0.4em] text-slate-400 dark:text-slate-500">
              Sua mensagem
            </label>
            <textarea
              id="chat-input"
              value={input}
              onChange={(event) => setInput(event.target.value)}
              placeholder={placeholder}
              rows={4}
              className="w-full resize-none rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200 dark:focus:border-slate-500 dark:focus:ring-slate-700/40"
            />
            <div className="mt-3 flex items-center justify-between">
              <p className="text-[11px] uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
                Histórico salvo localmente
              </p>
              <button
                type="submit"
                disabled={disabled}
                className="inline-flex items-center gap-2 rounded-full border border-slate-900 bg-slate-900 px-5 py-2 text-sm font-semibold text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] transition disabled:cursor-not-allowed disabled:opacity-60 dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]"
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
