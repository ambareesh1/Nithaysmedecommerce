"use client";

import { useState } from "react";
import { products as sampleProducts } from "@/data/products";
import { categories } from "@/data/products";

type AdminProduct = {
  id: number;
  name: string;
  category: string;
  price: number;
  stock: number;
  imageUrl: string;
  description: string;
};

const initialProducts: AdminProduct[] = sampleProducts.map((p) => ({
  id: p.id,
  name: p.name,
  category: p.category,
  price: p.price,
  stock: p.stock,
  imageUrl: "",
  description: p.description,
}));

export default function AdminProductsPage() {
  const [productList, setProductList] = useState<AdminProduct[]>(initialProducts);
  const [form, setForm] = useState({
    name: "",
    category: "Diagnostics",
    price: "",
    stock: "",
    imageUrl: "",
    description: "",
  });

  function updateForm(field: string, value: string) {
    setForm({ ...form, [field]: value });
  }

  function addProduct(e: React.FormEvent) {
    e.preventDefault();
    if (form.name === "" || form.price === "") {
      alert("Name and price are required");
      return;
    }
    const newProduct: AdminProduct = {
      id: Date.now(),
      name: form.name,
      category: form.category,
      price: Number(form.price),
      stock: Number(form.stock),
      imageUrl: form.imageUrl,
      description: form.description,
    };
    setProductList([newProduct, ...productList]);
    setForm({
      name: "",
      category: "Diagnostics",
      price: "",
      stock: "",
      imageUrl: "",
      description: "",
    });
  }

  function deleteProduct(id: number) {
    setProductList(productList.filter((p) => p.id !== id));
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-slate-800">Manage Products</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <form
          onSubmit={addProduct}
          className="bg-white border border-slate-200 rounded-lg p-6 flex flex-col gap-3"
        >
          <h2 className="font-semibold text-slate-800">Add Product</h2>

          <input
            type="text"
            placeholder="Product name"
            value={form.name}
            onChange={(e) => updateForm("name", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          />
          <select
            value={form.category}
            onChange={(e) => updateForm("category", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          >
            {categories.map((cat) => (
              <option key={cat} value={cat}>
                {cat}
              </option>
            ))}
          </select>
          <input
            type="number"
            placeholder="Price"
            value={form.price}
            onChange={(e) => updateForm("price", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          />
          <input
            type="number"
            placeholder="Stock"
            value={form.stock}
            onChange={(e) => updateForm("stock", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          />
          <input
            type="text"
            placeholder="Image URL"
            value={form.imageUrl}
            onChange={(e) => updateForm("imageUrl", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          />
          <textarea
            placeholder="Description"
            value={form.description}
            onChange={(e) => updateForm("description", e.target.value)}
            className="border border-slate-300 rounded-md px-3 py-2"
          />
          <button
            type="submit"
            className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
          >
            Add Product
          </button>
        </form>

        <div className="bg-white border border-slate-200 rounded-lg p-6 flex flex-col gap-2">
          <h2 className="font-semibold text-slate-800">Form Preview</h2>
          <p className="text-sm text-slate-600">Name: {form.name || "-"}</p>
          <p className="text-sm text-slate-600">
            Category: {form.category || "-"}
          </p>
          <p className="text-sm text-slate-600">Price: {form.price || "-"}</p>
          <p className="text-sm text-slate-600">Stock: {form.stock || "-"}</p>
          <p className="text-sm text-slate-600">
            Image URL: {form.imageUrl || "-"}
          </p>
          <p className="text-sm text-slate-600">
            Description: {form.description || "-"}
          </p>
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-lg p-6">
        <h2 className="font-semibold text-slate-800 mb-4">Product List</h2>
        {productList.length === 0 ? (
          <p className="text-center text-slate-500 py-6">
            No products added yet.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b border-slate-200 text-slate-600 text-sm">
                  <th className="py-2">Name</th>
                  <th className="py-2">Category</th>
                  <th className="py-2">Price</th>
                  <th className="py-2">Stock</th>
                  <th className="py-2">Action</th>
                </tr>
              </thead>
              <tbody>
                {productList.map((product) => (
                  <tr
                    key={product.id}
                    className="border-b border-slate-100 text-sm text-slate-700"
                  >
                    <td className="py-2">{product.name}</td>
                    <td className="py-2">{product.category}</td>
                    <td className="py-2">₹{product.price}</td>
                    <td className="py-2">{product.stock}</td>
                    <td className="py-2">
                      <button
                        onClick={() => deleteProduct(product.id)}
                        className="px-3 py-1 bg-red-100 text-red-700 rounded-md hover:bg-red-200"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
