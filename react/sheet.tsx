import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, invoker, pick } from "./util";

export type SheetProps = Omit<ComponentPropsWithRef<"dialog">, "title" | "id"> & {
  id: string;
  title: ReactNode;
  side?: "bottom" | "end";
  closeLabel?: string;
  closeButton?: boolean;
  dismissible?: boolean;
};

// Sheet is a modal <dialog> docked to the bottom, or with side="end" to the
// inline end. Open it like a Dialog: command="show-modal" commandfor={id}.
// dismissible adds closedby="any", so a tap on the backdrop closes it.
export function Sheet({ id, title, side, closeLabel = "Close", closeButton = true, dismissible, className, children, ...rest }: SheetProps) {
  const titleId = `${id}-title`;
  const closedBy: Record<string, string> = dismissible ? { closedby: "any" } : {};
  return (
    <dialog {...rest} {...closedBy} id={id} className={cx("ui-sheet", className)} aria-labelledby={titleId} data-side={pick(side, ["end"] as const)}>
      <header className="ui-sheet-head">
        <h2 className="ui-sheet-title" id={titleId}>{title}</h2>
        {closeButton && <button className="ui-sheet-x" type="button" aria-label={closeLabel} {...invoker("close", id)}>×</button>}
      </header>
      <div className="ui-sheet-body">{children}</div>
    </dialog>
  );
}
