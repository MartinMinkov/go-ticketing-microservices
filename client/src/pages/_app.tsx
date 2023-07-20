import "@/styles/globals.css";
import { AppContext, AppProps } from "next/app";
import { HTTPError } from "ky";

import { User, UserSchema } from "@/models/user";
import { buildClient } from "@/api/build-client";
import Header from "@/components/header";

type AppComponentProps = AppProps & {
  user: User | {};
};

export default function AppComponent({
  Component,
  pageProps,
  user,
}: AppComponentProps) {
  return (
    <div>
      <Header user={user} />
      <Component {...pageProps} />
    </div>
  );
}

AppComponent.getInitialProps = async (context: AppContext) => {
  const { req } = context.ctx;
  const data = await buildClient({
    host: req?.headers.host,
    cookie: req?.headers.cookie,
  })
    .get("api/users/currentuser")
    .catch((error: HTTPError) => {
      if (error.response?.status === 400) {
        return;
      }
      throw error;
    });

  if (!data) {
    return { user: {}, pageProps: {} };
  }
  const user = UserSchema.parse(await data.json());
  const pageProps = context.Component.getInitialProps
    ? await context.Component.getInitialProps(context.ctx)
    : {};
  return { user, pageProps };
};
