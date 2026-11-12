import type { ComponentPropsWithRef, ReactNode } from "react";
import { cx, pick } from "./util";

export type MenuProps = Omit<ComponentPropsWithRef<"div">, "id" | "title" | "popover" | "role"> & {
  id: string;
  label: string;
  sheet?: "phone";
  align?: "start" | "end";
  title?: ReactNode;
  closeLabel?: string;
};

// Menu is a popover anchored to the button that opens it: give that button
// popoverTarget={id} (and aria-haspopup="menu"). The browser closes it on
// Escape and an outside click and returns focus; ui.js adds the keys. With
// sheet="phone" it is a bottom sheet below 800px, with a head and a ×.
export function Menu({ id, label, sheet, align, title, closeLabel = "Close", className, children, ...rest }: MenuProps) {
  const asSheet = pick(sheet, ["phone"] as const);
  return (
    <div {...rest} id={id} className={cx("ui-menu", className)} popover="" role="menu" aria-label={label} data-align={pick(align, ["start"] as const)} data-sheet={asSheet}>
      {asSheet && (
        <div className="ui-menu-head" role="none">
          <span className="ui-menu-title">{title ?? label}</span>
          <button className="ui-menu-x" type="button" role="menuitem" popoverTarget={id} popoverTargetAction="hide" aria-label={closeLabel}>×</button>
        </div>
      )}
      {children}
    </div>
  );
}

type ItemOptions = {
  variant?: "danger";
  checked?: boolean;
  icon?: ReactNode;
  children?: ReactNode;
};

function itemAttrs({ variant, checked }: ItemOptions) {
  return {
    role: checked === undefined ? "menuitem" : "menuitemradio",
    "aria-checked": checked === undefined ? undefined : checked ? ("true" as const) : ("false" as const),
    "data-variant": pick(variant, ["danger"] as const),
  };
}

export type MenuItemProps = Omit<ComponentPropsWithRef<"button">, "type" | "role" | "children"> &
  ItemOptions & { type?: "button" | "submit" | "reset" };

// MenuItem is a row that does something. checked makes it a menuitemradio,
// which keeps the menu open when picked; any other item closes it (ui.js).
export function MenuItem({ variant, checked, icon, type, className, children, ...rest }: MenuItemProps) {
  return (
    <button {...rest} className={cx("ui-menu-item", className)} type={pick(type, ["submit", "reset"] as const) ?? "button"} {...itemAttrs({ variant, checked })}>
      {icon}
      {children}
    </button>
  );
}

export type MenuLinkProps = Omit<ComponentPropsWithRef<"a">, "role" | "children"> & ItemOptions;

// MenuLink is a row that goes somewhere.
export function MenuLink({ variant, checked, icon, className, children, ...rest }: MenuLinkProps) {
  return (
    <a {...rest} className={cx("ui-menu-item", className)} {...itemAttrs({ variant, checked })}>
      {icon}
      {children}
    </a>
  );
}

export type MenuLabelProps = Omit<ComponentPropsWithRef<"p">, "role">;

// MenuLabel is a caption or an identity row: not focusable, one line.
export function MenuLabel({ className, ...rest }: MenuLabelProps) {
  return <p {...rest} className={cx("ui-menu-label", className)} role="none" />;
}

export function MenuSeparator({ className, ...rest }: ComponentPropsWithRef<"hr">) {
  return <hr {...rest} className={cx("ui-menu-sep", className)} />;
}
