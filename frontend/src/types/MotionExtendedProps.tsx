// global.d.ts
import "framer-motion";
// Extensão da interface MotionProps para incluir className
declare module "framer-motion" {
  // Apenas estendemos MotionProps — removida a declaração de HTMLMotionProps que causava conflito.
  interface MotionProps extends React.HTMLAttributes<HTMLElement> {
    className?: string;
  }
}
// Preciso estenter a `disabled?: boolean;` para evitar erros em botões animados
declare global {
  interface HTMLMotionProps<T> {
    disabled?: boolean;
  }
}

export { };
