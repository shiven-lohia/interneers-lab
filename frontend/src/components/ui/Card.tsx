import React from "react";
import "./Card.css";

interface CardProps {
  children: React.ReactNode;
  onClick?: () => void;
  className?: string;
}

function Card({ children, onClick, className = "" }: CardProps) {
  const classes = ["card", onClick ? "card--clickable" : "", className]
    .filter(Boolean)
    .join(" ");

  return (
    <div className={classes} onClick={onClick}>
      {children}
    </div>
  );
}

export default Card;
