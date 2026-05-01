import React from "react";
import "./Button.css";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger";
}

function Button({
  variant = "primary",
  className = "",
  children,
  ...rest
}: ButtonProps) {
  return (
    <button className={`btn btn--${variant} ${className}`.trim()} {...rest}>
      {children}
    </button>
  );
}

export default Button;
