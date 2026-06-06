type AppButtonProps = {
  children: React.ReactNode;
  onClick?: () => void;
  type?: "button" | "submit";
  variant?: "primary" | "secondary";
};

export default function AppButton({
  children,
  onClick,
  type = "button",
  variant = "primary",
}: AppButtonProps) {
  const base = "px-4 py-2 rounded-md font-medium transition-colors";
  const styles =
    variant === "primary"
      ? "bg-blue-600 text-white hover:bg-blue-700"
      : "bg-slate-200 text-slate-800 hover:bg-slate-300";

  return (
    <button type={type} onClick={onClick} className={`${base} ${styles}`}>
      {children}
    </button>
  );
}
