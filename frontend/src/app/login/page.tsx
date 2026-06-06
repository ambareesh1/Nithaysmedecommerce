"use client";

import { useState } from "react";
import FormInput from "@/components/FormInput";
import AppButton from "@/components/AppButton";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (email === "" || password === "") {
      setError("Please fill in all fields");
      return;
    }
    if (!email.includes("@")) {
      setError("Please enter a valid email");
      return;
    }
    setError("");
    alert("Login successful for " + email);
  }

  return (
    <div className="max-w-md mx-auto bg-white border border-slate-200 rounded-lg p-6">
      <h1 className="text-2xl font-bold text-slate-800 mb-4">Login</h1>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <FormInput
          label="Email"
          type="email"
          value={email}
          placeholder="you@example.com"
          onChange={setEmail}
        />
        <FormInput
          label="Password"
          type="password"
          value={password}
          placeholder="Enter password"
          onChange={setPassword}
        />
        {error && <p className="text-sm text-red-600">{error}</p>}
        <AppButton type="submit">Login</AppButton>
      </form>
    </div>
  );
}
