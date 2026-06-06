export type CartItem = {
  id: number;
  name: string;
  category: string;
  price: number;
  quantity: number;
};

const CART_KEY = "medcart_items";

export function getCartItems(): CartItem[] {
  if (typeof window === "undefined") {
    return [];
  }
  const data = localStorage.getItem(CART_KEY);
  if (!data) {
    return [];
  }
  return JSON.parse(data);
}

export function saveCartItems(items: CartItem[]) {
  localStorage.setItem(CART_KEY, JSON.stringify(items));
}

export function addItemToCart(item: Omit<CartItem, "quantity">) {
  const items = getCartItems();
  const existing = items.find((i) => i.id === item.id);
  if (existing) {
    existing.quantity = existing.quantity + 1;
  } else {
    items.push({ ...item, quantity: 1 });
  }
  saveCartItems(items);
}

export function removeItemFromCart(id: number) {
  const items = getCartItems().filter((i) => i.id !== id);
  saveCartItems(items);
}

export function updateCartItemQuantity(id: number, quantity: number) {
  const items = getCartItems();
  const item = items.find((i) => i.id === id);
  if (item) {
    if (quantity < 1) {
      quantity = 1;
    }
    item.quantity = quantity;
  }
  saveCartItems(items);
}

export function clearCart() {
  saveCartItems([]);
}
