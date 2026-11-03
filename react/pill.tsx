import type { ComponentPropsWithRef } from "react";
import { cx, pick } from "./util";

export type PillVariant = "neutral" | "info" | "good" | "warning" | "bad";

const VARIANTS = ["info", "good", "warning", "bad"] as const;

export type PillProps = ComponentPropsWithRef<"span"> & { variant?: PillVariant };

export function Pill({ variant, className, ...rest }: PillProps) {
  return <span {...rest} className={cx("ui-pill", className)} data-variant={pick(variant, VARIANTS)} />;
}
