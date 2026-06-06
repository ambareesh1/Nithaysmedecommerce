"use client";

import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { products } from "@/data/products";
import { addItemToCart } from "@/lib/cart";

export default function ProductDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const id = Number(params.id);
  const product = products.find((p) => p.id === id);

  if (!product) {
    return (
      <div className="flex flex-col gap-4 items-start">
        <div className="bg-white border border-slate-200 rounded-lg p-10 text-center text-slate-500 w-full">
          Product not found.
        </div>
        <Link
          href="/products"
          className="px-4 py-2 rounded-md bg-slate-200 text-slate-800 hover:bg-slate-300"
        >
          Back to Products
        </Link>
      </div>
    );
  }

  function handleAddToCart() {
    addItemToCart({
      id: product!.id,
      name: product!.name,
      category: product!.category,
      price: product!.price,
    });
    alert(product!.name + " added to cart");
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="bg-white border border-slate-200 rounded-lg p-6 flex flex-col gap-3">
        <h1 className="text-2xl font-bold text-slate-800">{product.name}</h1>
        <span className="text-xs font-medium text-green-700 bg-green-100 px-2 py-1 rounded w-fit">
          {product.category}
        </span>
        <p className="text-2xl text-blue-600 font-bold">₹{product.price}</p>
        <p className="text-sm text-slate-500">In stock: {product.stock}</p>
        <p className="text-slate-700">{product.description}</p>

        <div className="flex gap-3 mt-3">
          <button
            onClick={handleAddToCart}
            className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
          >
            Add to Cart
          </button>
          <button
            onClick={() => router.push("/products")}
            className="px-4 py-2 rounded-md bg-slate-200 text-slate-800 hover:bg-slate-300"
          >
            Back to Products
          </button>
        </div>
      </div>
    </div>
  );
}
