import { CartItem } from "./cart";

export type OrderItem = {
  id: number;
  name: string;
  price: number;
  quantity: number;
  subtotal: number;
};

export type Order = {
  orderNumber: string;
  date: string;
  status: string;
  items: OrderItem[];
  totalQuantity: number;
  totalAmount: number;
};

const ORDERS_KEY = "medcart_orders";

export function getOrders(): Order[] {
  if (typeof window === "undefined") {
    return [];
  }
  const data = localStorage.getItem(ORDERS_KEY);
  if (!data) {
    return [];
  }
  return JSON.parse(data);
}

export function saveOrders(orders: Order[]) {
  localStorage.setItem(ORDERS_KEY, JSON.stringify(orders));
}

export function createOrderFromCart(cartItems: CartItem[]): Order {
  const items: OrderItem[] = cartItems.map((item) => ({
    id: item.id,
    name: item.name,
    price: item.price,
    quantity: item.quantity,
    subtotal: item.price * item.quantity,
  }));

  const totalQuantity = items.reduce((sum, item) => sum + item.quantity, 0);
  const totalAmount = items.reduce((sum, item) => sum + item.subtotal, 0);

  const order: Order = {
    orderNumber: "ORD-" + Date.now(),
    date: new Date().toLocaleString(),
    status: "Pending",
    items,
    totalQuantity,
    totalAmount,
  };

  const orders = getOrders();
  orders.unshift(order);
  saveOrders(orders);

  return order;
}
