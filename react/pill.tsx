import type { ComponentPropsWithRef } from "react";
import { cx, pick } from "./util";

export type PillVariant = "neutral" | "info" | "good" | "warning" | "bad";
export type PillTone = "fill" | "outline";

const VARIANTS = ["info", "good", "warning", "bad"] as const;
const TONES = ["outline"] as const;

export type PillProps = ComponentPropsWithRef<"span"> & { variant?: PillVariant; tone?: PillTone; dot?: boolean };

export function Pill({ variant, tone, dot, className, ...rest }: PillProps) {
  return (
    <span
      {...rest}
      className={cx("ui-pill", className)}
      data-variant={pick(variant, VARIANTS)}
      data-tone={pick(tone, TONES)}
      data-dot={dot ? "" : undefined}
    />
  );
}
