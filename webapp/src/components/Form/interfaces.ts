import { ChangeEventHandler } from "react";

export type InputChangeHandler = ChangeEventHandler<HTMLInputElement>;
export type FormDataHandler = (form: Record<string, string>) => void;
