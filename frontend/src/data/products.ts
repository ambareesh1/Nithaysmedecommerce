export type Product = {
  id: number;
  name: string;
  category: string;
  price: number;
  stock: number;
  description: string;
};

export const products: Product[] = [
  {
    id: 1,
    name: "Digital Thermometer",
    category: "Diagnostics",
    price: 499,
    stock: 25,
    description: "A digital thermometer used to measure body temperature quickly and accurately.",
  },
  {
    id: 2,
    name: "Surgical Gloves",
    category: "Protection",
    price: 299,
    stock: 100,
    description: "Disposable latex surgical gloves for safe medical handling.",
  },
  {
    id: 3,
    name: "Blood Pressure Monitor",
    category: "Diagnostics",
    price: 1999,
    stock: 15,
    description: "Automatic blood pressure monitor suitable for home and clinic use.",
  },
  {
    id: 4,
    name: "Pulse Oximeter",
    category: "Diagnostics",
    price: 899,
    stock: 40,
    description: "Fingertip pulse oximeter to measure blood oxygen levels and pulse rate.",
  },
  {
    id: 5,
    name: "Nebulizer Machine",
    category: "Equipment",
    price: 2499,
    stock: 10,
    description: "Compact nebulizer machine for effective respiratory care at home.",
  },
  {
    id: 6,
    name: "First Aid Kit",
    category: "Supplies",
    price: 799,
    stock: 50,
    description: "Complete first aid kit with essential items for emergency care.",
  },
];

export const categories = ["Diagnostics", "Protection", "Equipment", "Supplies"];
