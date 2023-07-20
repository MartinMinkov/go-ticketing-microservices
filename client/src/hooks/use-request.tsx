import { useState } from "react";
import ky, { HTTPError } from "ky";
import { z } from "zod";

type UseRequestInput = {
  url: string;
  method: "get" | "post" | "put" | "delete";
};

export default function useRequest({ url, method }: UseRequestInput) {
  const [requestErrors, setRequestErrors] = useState<React.ReactNode>(null);
  const doRequest = async <T,>(
    body: unknown,
    schema: z.ZodObject<T extends z.ZodRawShape ? T : never> | undefined
  ) => {
    try {
      setRequestErrors(null);
      const response = await ky[method](url, {
        json: body,
      }).json();
      console.log("DEBUG HERE3", response);
      return schema ? schema.parse(response) : response;
    } catch (err) {
      if (err instanceof HTTPError) {
        const errorJson = await err.response.json();
        setRequestErrors(
          <div
            role="alert"
            className="relative rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700"
          >
            <strong className="font-bold">Oops...</strong>
            <ul role="list" className="list-disc space-y-1 pl-5">
              {errorJson.errors.map((err: string) => {
                return <li key={err}>{err}</li>;
              })}
            </ul>
          </div>
        );
      }
    }
  };

  return { doRequest, requestErrors };
}
