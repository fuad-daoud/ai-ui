import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, pick } from "./util";

const TYPES = ["radio"] as const;

export type CheckProps = Omit<ComponentPropsWithRef<"input">, "type" | "className"> & {
  label: ReactNode;
  type?: "checkbox" | "radio";
  className?: string;
};

export function Check({ label, type, className, ...rest }: CheckProps) {
  return (
    <label className={cx("ui-check", className)}>
      <input {...rest} className="ui-check-input" type={pick(type, TYPES) ?? "checkbox"} />
      <span className="ui-check-label">{label}</span>
    </label>
  );
}
