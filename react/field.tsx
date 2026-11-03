import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, pick, present } from "./util";

const INPUT_TYPES = ["email", "password", "number", "search", "tel", "url", "date"] as const;
export type InputType = "text" | (typeof INPUT_TYPES)[number];
export type FieldOption = { value: string; label: string };

type FieldBase = {
  label: ReactNode;
  name?: string;
  id?: string;
  hint?: ReactNode;
  error?: ReactNode;
  className?: string;
};
type Own = keyof FieldBase | "control";

export type InputFieldProps = FieldBase & { control?: "input"; type?: InputType } & Omit<ComponentPropsWithRef<"input">, Own | "type">;
export type TextareaFieldProps = FieldBase & { control: "textarea" } & Omit<ComponentPropsWithRef<"textarea">, Own>;
export type SelectFieldProps = FieldBase & { control: "select"; options: readonly FieldOption[] } & Omit<ComponentPropsWithRef<"select">, Own>;
export type FieldProps = InputFieldProps | TextareaFieldProps | SelectFieldProps;

type Wiring = { id: string; name?: string; "aria-describedby"?: string; "aria-invalid"?: "true" };

export function Field(props: FieldProps) {
  const { label, name, id: idProp, hint, error, className, ...controlProps } = props;
  const id = idProp || name;
  if (!id) throw new Error("ui.field: needs ID or Name");
  const hasHint = present(hint);
  const hasError = present(error);
  const describedBy = [hasHint ? `${id}-hint` : "", hasError ? `${id}-error` : ""].filter(Boolean).join(" ");
  const wiring: Wiring = {
    id,
    name,
    "aria-describedby": describedBy || undefined,
    "aria-invalid": hasError ? "true" : undefined,
  };
  return (
    <div className={cx("ui-field", className)}>
      <label className="ui-field-label" htmlFor={id}>{label}</label>
      {control(controlProps, wiring)}
      {hasHint && <p className="ui-field-hint" id={`${id}-hint`}>{hint}</p>}
      {hasError && <p className="ui-field-error" id={`${id}-error`}>{error}</p>}
    </div>
  );
}

type ControlProps = {
  control?: string;
  type?: string;
  options?: readonly FieldOption[];
  defaultValue?: string | number | readonly string[];
};

// blank drops an empty defaultValue, so an empty input renders no value
// attribute, matching the Go partial.
const blank = (v: ControlProps["defaultValue"]) => (v === "" ? undefined : v);

function control(props: object, wiring: Wiring): ReactNode {
  const { control: kind, type, options, defaultValue, ...rest } = props as ControlProps;
  if (kind === "textarea") {
    return <textarea {...(rest as ComponentPropsWithRef<"textarea">)} {...wiring} className="ui-textarea" defaultValue={blank(defaultValue)} />;
  }
  if (kind === "select") {
    return (
      <select {...(rest as ComponentPropsWithRef<"select">)} {...wiring} className="ui-select" defaultValue={defaultValue}>
        {(options ?? []).map((o) => (
          <option key={o.value} value={o.value}>{o.label}</option>
        ))}
      </select>
    );
  }
  return <input {...(rest as ComponentPropsWithRef<"input">)} {...wiring} className="ui-input" type={pick(type, INPUT_TYPES) ?? "text"} defaultValue={blank(defaultValue)} />;
}
