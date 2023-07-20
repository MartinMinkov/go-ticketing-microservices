import SignUpForm from "@/components/sign-up-form";

const SignIn = () => {
  return (
    <main className="container px-8">
      <h1 className="text-xl font-bold tracking-tight text-gray-900 sm:text-4xl">
        Sign In
      </h1>
      <div className="pt-6 pb-8">
        <SignUpForm apiUrl={"/api/users/signin"} buttonText={"Sign In"} />
      </div>
    </main>
  );
};
export default SignIn;
