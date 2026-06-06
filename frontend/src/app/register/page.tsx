"use client";

import { useState } from "react";
import FormInput from "@/components/FormInput";
import AppButton from "@/components/AppButton";

export default function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (name === "" || email === "" || password === "") {
      setError("Please fill in all fields");
      return;
    }
    if (!email.includes("@")) {
      setError("Please enter a valid email");
      return;
    }
    if (password.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }
    setError("");
    alert("Registration successful for " + name);
  }

  return (
    <div className="max-w-md mx-auto bg-white border border-slate-200 rounded-lg p-6">
      <h1 className="text-2xl font-bold text-slate-800 mb-4">Register</h1>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <FormInput
          label="Name"
          value={name}
          placeholder="Your name"
          onChange={setName}
        />
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
        <AppButton type="submit">Register</AppButton>
      </form>
    </div>
  );
}
