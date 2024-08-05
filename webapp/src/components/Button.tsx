import React, { ButtonHTMLAttributes } from 'react'
import { twMerge } from 'tailwind-merge'

const Button: React.FC<ButtonHTMLAttributes<HTMLButtonElement>> = ({ className, ...properties }) => {
  const classes = twMerge(
    'w-full h-[40px] bg-primary rounded-lg transition text-[#FFF] font-bold',
    className
  )

  return <button {...properties} className={classes} />
}

export default Button
