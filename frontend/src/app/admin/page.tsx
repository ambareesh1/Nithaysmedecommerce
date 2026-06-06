"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { products } from "@/data/products";
import { getOrders } from "@/lib/orders";

export default function AdminPage() {
  const [totalOrders, setTotalOrders] = useState(0);

  useEffect(() => {
    setTotalOrders(getOrders().length);
  }, []);

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-slate-800">Admin Dashboard</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Link
          href="/admin/products"
          className="bg-white border border-slate-200 rounded-lg p-6 hover:border-blue-400"
        >
          <h3 className="font-semibold text-blue-600">Manage Products</h3>
          <p className="text-sm text-slate-600 mt-2">
            Add, view, and delete medical products.
          </p>
        </Link>

        <Link
          href="/orders"
          className="bg-white border border-slate-200 rounded-lg p-6 hover:border-blue-400"
        >
          <h3 className="font-semibold text-blue-600">View Orders</h3>
          <p className="text-sm text-slate-600 mt-2">
            See all placed orders and their status.
          </p>
        </Link>

        <div className="bg-white border border-slate-200 rounded-lg p-6">
          <h3 className="font-semibold text-slate-700">Total Products</h3>
          <p className="text-3xl font-bold text-blue-600 mt-2">
            {products.length}
          </p>
        </div>

        <div className="bg-white border border-slate-200 rounded-lg p-6">
          <h3 className="font-semibold text-slate-700">Total Orders</h3>
          <p className="text-3xl font-bold text-blue-600 mt-2">
            {totalOrders}
          </p>
        </div>
      </div>
    </div>
  );
}
