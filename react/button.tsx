import type { ComponentPropsWithRef } from "react";
import { cx, invoker, pick } from "./util";

export type ButtonVariant = "primary" | "secondary" | "ghost" | "danger";
export type ButtonSize = "sm" | "md" | "lg";

const VARIANTS = ["secondary", "ghost", "danger"] as const;
const SIZES = ["sm", "lg"] as const;
const TYPES = ["submit", "reset"] as const;

export type ButtonProps = Omit<ComponentPropsWithRef<"button">, "type"> & {
  variant?: ButtonVariant;
  size?: ButtonSize;
  type?: "button" | "submit" | "reset";
  command?: string;
  commandFor?: string;
};

export function Button({ variant, size, type, command, commandFor, className, ...rest }: ButtonProps) {
  return (
    <button
      {...rest}
      {...invoker(command, commandFor)}
      className={cx("ui-btn", className)}
      type={pick(type, TYPES) ?? "button"}
      data-variant={pick(variant, VARIANTS)}
      data-size={pick(size, SIZES)}
    />
  );
}

export type ButtonLinkProps = ComponentPropsWithRef<"a"> & {
  variant?: ButtonVariant;
  size?: ButtonSize;
  disabled?: boolean;
};

export function ButtonLink({ variant, size, disabled, href, className, ...rest }: ButtonLinkProps) {
  return (
    <a
      {...rest}
      className={cx("ui-btn", className)}
      href={disabled ? undefined : href}
      aria-disabled={disabled ? "true" : undefined}
      data-variant={pick(variant, VARIANTS)}
      data-size={pick(size, SIZES)}
    />
  );
}
