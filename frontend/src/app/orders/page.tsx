"use client";

import { useEffect, useState } from "react";
import { Order, getOrders } from "@/lib/orders";

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    setOrders(getOrders());
  }, []);

  if (orders.length === 0) {
    return (
      <div className="bg-white border border-slate-200 rounded-lg p-10 text-center text-slate-500">
        You have no orders yet. Place an order to see it here.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-slate-800">Your Orders</h1>

      {orders.map((order) => (
        <div
          key={order.orderNumber}
          className="bg-white border border-slate-200 rounded-lg p-4 flex flex-col gap-3"
        >
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <h3 className="font-semibold text-slate-800">
                {order.orderNumber}
              </h3>
              <p className="text-sm text-slate-500">{order.date}</p>
            </div>
            <span className="text-xs font-medium text-yellow-700 bg-yellow-100 px-3 py-1 rounded w-fit">
              {order.status}
            </span>
          </div>

          <div className="border-t border-slate-100 pt-3 flex flex-col gap-2">
            {order.items.map((item) => (
              <div
                key={item.id}
                className="flex justify-between text-sm text-slate-700"
              >
                <span>
                  {item.name} x {item.quantity}
                </span>
                <span>₹{item.subtotal}</span>
              </div>
            ))}
          </div>

          <div className="border-t border-slate-100 pt-3 flex justify-between">
            <span className="text-slate-700">
              Total Quantity: {order.totalQuantity}
            </span>
            <span className="font-bold text-slate-800">
              Total: ₹{order.totalAmount}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}
