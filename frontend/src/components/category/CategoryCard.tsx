import "./CategoryCard.css";
import type { Category } from "../../types";

interface CategoryCardProps {
  category: Category;
  onClick: () => void;
}

function CategoryCard({ category, onClick }: CategoryCardProps) {
  const hasDescription =
    category.description && category.description.trim() !== "";
  return (
    <div
      className="category-card"
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => (e.key === "Enter" || e.key === " ") && onClick()}
    >
      <span className="category-card__title">{category.title}</span>
      <span className="category-card__description">
        {hasDescription ? category.description : "No description"}
      </span>
    </div>
  );
}

export default CategoryCard;
