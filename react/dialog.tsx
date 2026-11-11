import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, invoker, present } from "./util";

export type DialogProps = Omit<ComponentPropsWithRef<"dialog">, "title" | "id"> & {
  id: string;
  title: ReactNode;
  closeLabel?: string;
  closeButton?: boolean;
  eyebrow?: ReactNode;
};

// Dialog is a modal <dialog>. Open it from any button with
// command="show-modal" commandfor={id}; the × button closes it with
// command="close". It needs no JavaScript (spec §4.3). Without the ×
// (closeButton={false}), give it a button that closes it.
export function Dialog({ id, title, closeLabel = "Close", closeButton = true, eyebrow, className, children, ...rest }: DialogProps) {
  const titleId = `${id}-title`;
  return (
    <dialog {...rest} id={id} className={cx("ui-dialog", className)} aria-labelledby={titleId}>
      {present(eyebrow) && <p className="ui-dialog-eyebrow">{eyebrow}</p>}
      <header className="ui-dialog-head">
        <h2 className="ui-dialog-title" id={titleId}>{title}</h2>
        {closeButton && <button className="ui-dialog-x" type="button" aria-label={closeLabel} {...invoker("close", id)}>×</button>}
      </header>
      <div className="ui-dialog-body">{children}</div>
    </dialog>
  );
}
