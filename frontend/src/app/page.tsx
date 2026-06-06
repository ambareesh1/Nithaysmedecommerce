import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col items-center text-center gap-6 py-12">
      <h1 className="text-4xl font-bold text-slate-800">MedCart SaaS</h1>
      <p className="max-w-2xl text-slate-600 text-lg">
        A medical product ordering platform for hospitals, clinics, and
        diagnostic centers. Browse medical products, add them to your cart, and
        place orders easily.
      </p>
      <div className="flex gap-4 flex-wrap justify-center">
        <Link
          href="/products"
          className="px-6 py-3 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700"
        >
          Explore Products
        </Link>
        <Link
          href="/admin"
          className="px-6 py-3 rounded-md bg-slate-200 text-slate-800 font-medium hover:bg-slate-300"
        >
          Admin Dashboard
        </Link>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-10 w-full">
        <div className="bg-white border border-slate-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-600">Wide Product Range</h3>
          <p className="text-sm text-slate-600 mt-2">
            Thermometers, monitors, gloves, and more medical supplies.
          </p>
        </div>
        <div className="bg-white border border-slate-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-600">Easy Ordering</h3>
          <p className="text-sm text-slate-600 mt-2">
            Add to cart and place orders in just a few clicks.
          </p>
        </div>
        <div className="bg-white border border-slate-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-600">Order History</h3>
          <p className="text-sm text-slate-600 mt-2">
            Track all your past orders in one place.
          </p>
        </div>
      </div>
    </div>
  );
}
