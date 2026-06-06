"use client";

import Link from "next/link";
import { addItemToCart } from "@/lib/cart";

type ProductCardProps = {
  id: number;
  name: string;
  category: string;
  price: number;
  stock: number;
};

export default function ProductCard({
  id,
  name,
  category,
  price,
  stock,
}: ProductCardProps) {
  function handleAddToCart() {
    addItemToCart({ id, name, category, price });
    alert(name + " added to cart");
  }

  return (
    <div className="bg-white border border-slate-200 rounded-lg p-4 flex flex-col gap-2 shadow-sm">
      <h3 className="text-lg font-semibold text-slate-800">{name}</h3>
      <span className="text-xs font-medium text-green-700 bg-green-100 px-2 py-1 rounded w-fit">
        {category}
      </span>
      <p className="text-blue-600 font-bold">₹{price}</p>
      <p className="text-sm text-slate-500">In stock: {stock}</p>
      <div className="flex gap-2 mt-2">
        <Link
          href={`/products/${id}`}
          className="px-3 py-1 rounded-md bg-slate-200 text-slate-800 hover:bg-slate-300 text-sm"
        >
          View
        </Link>
        <button
          onClick={handleAddToCart}
          className="px-3 py-1 rounded-md bg-blue-600 text-white hover:bg-blue-700 text-sm"
        >
          Add to Cart
        </button>
      </div>
    </div>
  );
}
