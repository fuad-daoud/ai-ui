"use client";
import { useEffect, type ComponentPropsWithRef } from "react";
import { cx } from "./util";

export type SwatchesProps = Omit<ComponentPropsWithRef<"div">, "role"> & { label: string };

// Swatches is a row of theme choices.
export function Swatches({ label, className, ...rest }: SwatchesProps) {
  return <div {...rest} className={cx("ui-swatches", className)} role="group" aria-label={label} />;
}

export type SwatchTheme = "light" | "dark" | "system";

export type SwatchProps = Omit<ComponentPropsWithRef<"button">, "children" | "type"> & {
  theme: SwatchTheme;
  palette: string;
  label: string;
};

type ThemeWindow = { aiui?: { theme?: { sync?: () => void } } };

// Swatch is a button that picks a theme and palette through theme.js, which
// also marks the current one aria-pressed. Its face previews the choice.
export function Swatch({ theme, palette, label, className, ...rest }: SwatchProps) {
  useEffect(() => {
    (window as unknown as ThemeWindow).aiui?.theme?.sync?.();
  }, [theme, palette]);
  return (
    <button {...rest} className={cx("ui-swatch", className)} type="button" data-ui-set-theme={theme} data-ui-set-palette={palette} aria-label={label}>
      <SwatchFace theme={theme} palette={palette} />
    </button>
  );
}

export type SwatchFaceProps = Omit<ComponentPropsWithRef<"span">, "children"> & { theme: SwatchTheme; palette: string };

// SwatchFace previews a palette and mode; it can also lead a MenuItem.
export function SwatchFace({ theme, palette, className, ...rest }: SwatchFaceProps) {
  return <span {...rest} className={cx("ui-swatch-face", className)} data-theme={theme} data-palette={palette} aria-hidden="true" />;
}
