'use client'

import Button from "@/components/Button";
import { Form, FormDataHandler } from "@/components/Form";
import { apiService } from "@/services/api";
import React from "react";

const LoginForm: React.FC = () => {
  const formDataHandler: FormDataHandler = (data) => {
    apiService
      .login({ email: data.email, password: data.password })
      .then(console.log)
      .catch(console.error)
  }

  return (
    <Form.Container
      formData={formDataHandler}
      className="w-full bg-cards rounded-xl mx-auto max-w-[400px] p-6 flex flex-col gap-4"
    >
      <Form.Input label="E-mail" id="email" type="email" required />
      <Form.Input label="Password" id="password" type="password" minLength={8} required />

      <Button type="submit" className="mt-2">Entrar</Button>
    </Form.Container>
  )
}

export default LoginForm
