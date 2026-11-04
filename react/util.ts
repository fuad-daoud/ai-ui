export function cx(...parts: (string | false | null | undefined)[]): string {
  return parts.filter(Boolean).join(" ");
}

// pick returns value if it is one of allowed, else undefined. React then
// omits the attribute and the component's default applies (spec §6).
export function pick<T extends string>(value: unknown, allowed: readonly T[]): T | undefined {
  return (allowed as readonly unknown[]).includes(value) ? (value as T) : undefined;
}

export function present(v: unknown): boolean {
  return v !== undefined && v !== null && v !== false && v !== "";
}

// invoker returns the lowercase command / commandfor attributes of an
// invoker button (spec §4.3), omitting the empty ones.
export function invoker(command?: string, commandFor?: string): Record<string, string> {
  const a: Record<string, string> = {};
  if (command) a.command = command;
  if (commandFor) a.commandfor = commandFor;
  return a;
}
