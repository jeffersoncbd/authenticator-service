'use client'

import React, { InputHTMLAttributes, useContext } from 'react'
import { twMerge } from 'tailwind-merge'
import { FormContext } from './Context'

interface InputProperties extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  id: string
}


const FormInput: React.FC<InputProperties> = ({ label, className, ...properties }) => {
  const formContext = useContext(FormContext)
  const classes = twMerge([
    'rounded-lg px-2 h-[40px] w-full bg-gray-200 text-black',
    className
  ])

  return (
    <div>
      <label htmlFor={properties.id}>{label}:</label>
      <input
        onChange={formContext.inputChangeHandler}
        {...properties}
        className={classes}
      />
    </div>
  )
}

export default FormInput
