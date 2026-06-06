import Link from "next/link";

const links = [
  { name: "Home", href: "/" },
  { name: "Products", href: "/products" },
  { name: "Cart", href: "/cart" },
  { name: "Orders", href: "/orders" },
  { name: "Admin", href: "/admin" },
  { name: "Login", href: "/login" },
];

export default function Header() {
  return (
    <header className="bg-white border-b border-slate-200 shadow-sm">
      <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="text-xl font-bold text-blue-600">
          MedCart SaaS
        </Link>
        <nav className="flex gap-4 flex-wrap">
          {links.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="text-slate-700 hover:text-blue-600 font-medium"
            >
              {link.name}
            </Link>
          ))}
        </nav>
      </div>
    </header>
  );
}
