import SignUpForm from "@/components/sign-up-form";

const SignUp = () => {
  return (
    <main className="container px-8">
      <h1 className="text-xl font-bold tracking-tight text-gray-900 sm:text-4xl">
        Sign Up
      </h1>
      <div className="pt-6 pb-8">
        <SignUpForm apiUrl={"/api/users/signup"} buttonText={"Sign Up"} />
      </div>
    </main>
  );
};

export default SignUp;
