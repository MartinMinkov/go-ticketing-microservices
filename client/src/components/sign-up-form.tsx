import Router from "next/router";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import useRequest from "@/hooks/use-request";

const formSchema = z.object({
  email: z.string().email({ message: "Please enter a valid email address" }),
  password: z
    .string()
    .min(6, { message: "Password must be at least 6 characters" }),
});
type FormValues = z.infer<typeof formSchema>;

type SignUpFormProps = {
  apiUrl: string;
  buttonText: string;
};

const SignUpForm = ({ apiUrl, buttonText }: SignUpFormProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(formSchema) });

  const { doRequest, requestErrors } = useRequest({
    url: apiUrl,
    method: "post",
  });

  const onSubmit = async (formData: FormValues) => {
    const { email, password } = formData;
    console.log("DBEUG HERE", email, password);
    const user = await doRequest(
      {
        email,
        password,
      },
      z.object({
        id: z.string(),
        email: z.string().email(),
        created_at: z.string(),
        updated_at: z.string(),
      })
    );
    console.log("DBEUG HERE2", email, password, user);
    if (user) Router.push("/");
  };

  return (
    <form className="rounded bg-white" onSubmit={handleSubmit(onSubmit)}>
      <div className="mb-4">
        <label
          className="mb-2 block text-sm font-bold text-gray-700"
          htmlFor="username"
        >
          Email Address
        </label>
        <input
          {...register("email", { required: true })}
          className={`focus:shadow-outline mb-3 w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none ${
            errors.email ? "border-red-500" : ""
          }`}
          type="email"
          placeholder="Email Address"
        />
        {errors.email?.message && (
          <p className="text-xs italic text-red-500">
            {errors.email?.message.toString()}
          </p>
        )}
      </div>
      <div className="mb-6">
        <label
          className="mb-2 block text-sm font-bold text-gray-700"
          htmlFor="password"
        >
          Password
        </label>
        <input
          {...register("password", {
            required: true,
            minLength: 6,
          })}
          className={`focus:shadow-outline mb-3 w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none ${
            errors.password ? "border-red-500" : ""
          }`}
          type="password"
        />
        {errors.password?.message && (
          <p className="text-xs italic text-red-500">
            {errors.password?.message.toString()}
          </p>
        )}
      </div>
      <div className="flex items-center justify-between">
        <button
          className="focus:shadow-outline rounded bg-blue-500 py-2 px-4 font-bold text-white hover:bg-blue-700 focus:outline-none"
          type="submit"
        >
          {buttonText}
        </button>
      </div>
      {requestErrors ? requestErrors : null}
    </form>
  );
};

export default SignUpForm;
