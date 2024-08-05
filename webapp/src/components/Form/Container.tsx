'use client'

import React, { FormHTMLAttributes, useState } from "react";
import { FormContext } from "./Context";
import { FormDataHandler, InputChangeHandler } from "./interfaces";


interface FormContainerProperties extends FormHTMLAttributes<HTMLFormElement> {
  formData: FormDataHandler
}

const FormContainer: React.FC<FormContainerProperties> = ({ formData, ...properties }) => {
  const [form, setForm] = useState({})

  const inputChangeHandler: InputChangeHandler = (event) => {
    setForm({ ...form, [event.target.id]: event.target.value })
  }

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    formData(form)
  }

  return (
    <FormContext.Provider value={{ inputChangeHandler }}>
      <form {...properties} onSubmit={handleSubmit} />
    </FormContext.Provider>
  )
}

export default FormContainer
