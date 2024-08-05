import { createContext } from "react";
import { InputChangeHandler } from "./interfaces";

interface FormContextProperties {
  inputChangeHandler: InputChangeHandler;
}

export const FormContext = createContext<FormContextProperties>({
  inputChangeHandler: () => undefined,
});
