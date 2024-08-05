import LoginForm from "./form";

export default function Login() {
  return (
    <div className="h-[100vh] flex flex-col justify-center items-center px-4 gap-6">
      <p className="text-6xl">🔐</p>
      <h1 className="text-2xl font-bold">Authenticator</h1>

      <LoginForm />
    </div>
  )
}
