import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, invoker, pick } from "./util";

export type ButtonVariant = "primary" | "secondary" | "ghost" | "danger";
export type ButtonSize = "sm" | "md" | "lg";

const VARIANTS = ["secondary", "ghost", "danger"] as const;
const SIZES = ["sm", "lg"] as const;
const TYPES = ["submit", "reset"] as const;
const TONES = ["outline"] as const;
const SHAPES = ["square"] as const;
const WIDTHS = ["full"] as const;

export type ButtonTone = "fill" | "outline";
export type ButtonShape = "auto" | "square";
export type ButtonWidth = "auto" | "full";

type Options = {
  variant?: ButtonVariant;
  size?: ButtonSize;
  tone?: ButtonTone;
  shape?: ButtonShape;
  width?: ButtonWidth;
  icon?: ReactNode;
};

// options returns the data attributes and the content shared by Button and
// ButtonLink. A square button shows only its icon, and its text label becomes
// the accessible name, as in the Go partial.
function options({ variant, size, tone, shape, width, icon }: Options, children: ReactNode) {
  const square = pick(shape, SHAPES);
  if (square && typeof children !== "string") throw new Error("ui.button: a square button's label must be text");
  return {
    attrs: {
      "data-variant": pick(variant, VARIANTS),
      "data-size": pick(size, SIZES),
      "data-tone": pick(tone, TONES),
      "data-shape": square,
      "aria-label": square ? (children as string) : undefined,
      "data-width": pick(width, WIDTHS),
    },
    content: square ? icon : <>{icon}{children}</>,
  };
}

export type ButtonProps = Omit<ComponentPropsWithRef<"button">, "type" | "children"> & Options & {
  children?: ReactNode;
  type?: "button" | "submit" | "reset";
  command?: string;
  commandFor?: string;
};

export function Button({ variant, size, tone, shape, width, icon, type, command, commandFor, className, children, ...rest }: ButtonProps) {
  const o = options({ variant, size, tone, shape, width, icon }, children);
  return (
    <button {...rest} {...invoker(command, commandFor)} className={cx("ui-btn", className)} type={pick(type, TYPES) ?? "button"} {...o.attrs}>
      {o.content}
    </button>
  );
}

export type ButtonLinkProps = Omit<ComponentPropsWithRef<"a">, "children"> & Options & {
  children?: ReactNode;
  disabled?: boolean;
};

export function ButtonLink({ variant, size, tone, shape, width, icon, disabled, href, className, children, ...rest }: ButtonLinkProps) {
  const o = options({ variant, size, tone, shape, width, icon }, children);
  return (
    <a {...rest} className={cx("ui-btn", className)} href={disabled ? undefined : href} aria-disabled={disabled ? "true" : undefined} {...o.attrs}>
      {o.content}
    </a>
  );
}
