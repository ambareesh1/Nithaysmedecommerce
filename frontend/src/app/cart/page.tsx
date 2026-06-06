"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  CartItem,
  getCartItems,
  removeItemFromCart,
  updateCartItemQuantity,
  clearCart,
} from "@/lib/cart";
import { createOrderFromCart } from "@/lib/orders";

export default function CartPage() {
  const router = useRouter();
  const [items, setItems] = useState<CartItem[]>([]);

  useEffect(() => {
    setItems(getCartItems());
  }, []);

  function refresh() {
    setItems(getCartItems());
  }

  function increase(id: number, quantity: number) {
    updateCartItemQuantity(id, quantity + 1);
    refresh();
  }

  function decrease(id: number, quantity: number) {
    updateCartItemQuantity(id, quantity - 1);
    refresh();
  }

  function remove(id: number) {
    removeItemFromCart(id);
    refresh();
  }

  function handleClear() {
    clearCart();
    refresh();
  }

  function placeOrder() {
    createOrderFromCart(items);
    clearCart();
    router.push("/orders");
  }

  const totalQuantity = items.reduce((sum, item) => sum + item.quantity, 0);
  const totalAmount = items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0
  );

  if (items.length === 0) {
    return (
      <div className="bg-white border border-slate-200 rounded-lg p-10 text-center text-slate-500">
        Your cart is empty. Add some products to get started.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-slate-800">Your Cart</h1>

      <div className="flex flex-col gap-3">
        {items.map((item) => (
          <div
            key={item.id}
            className="bg-white border border-slate-200 rounded-lg p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
          >
            <div>
              <h3 className="font-semibold text-slate-800">{item.name}</h3>
              <p className="text-sm text-slate-500">{item.category}</p>
              <p className="text-blue-600 font-bold">₹{item.price}</p>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => decrease(item.id, item.quantity)}
                className="px-3 py-1 bg-slate-200 rounded-md hover:bg-slate-300"
              >
                -
              </button>
              <span className="w-8 text-center">{item.quantity}</span>
              <button
                onClick={() => increase(item.id, item.quantity)}
                className="px-3 py-1 bg-slate-200 rounded-md hover:bg-slate-300"
              >
                +
              </button>
              <button
                onClick={() => remove(item.id)}
                className="px-3 py-1 bg-red-100 text-red-700 rounded-md hover:bg-red-200"
              >
                Remove
              </button>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-white border border-slate-200 rounded-lg p-4 flex flex-col gap-2">
        <p className="text-slate-700">Total Quantity: {totalQuantity}</p>
        <p className="text-lg font-bold text-slate-800">
          Total Amount: ₹{totalAmount}
        </p>
        <div className="flex gap-3 mt-2">
          <button
            onClick={placeOrder}
            className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
          >
            Place Order
          </button>
          <button
            onClick={handleClear}
            className="px-4 py-2 rounded-md bg-slate-200 text-slate-800 hover:bg-slate-300"
          >
            Clear Cart
          </button>
        </div>
      </div>
    </div>
  );
}
