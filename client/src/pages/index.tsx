import Head from "next/head";
import { NextPage } from "next";
import { HTTPError } from "ky";

import { User, UserSchema, isUserEmpty } from "@/models/user";
import { buildClient } from "@/api/build-client";

const Home: NextPage<User | {}> = (user: User | {}) => {
  return (
    <>
      <Head>
        <title>Git Tix Ticketing</title>
        <meta name="description" content="Ticketing Microservices" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className="container mx-auto">
        <div>Boom goes the dynamite</div>
        <div>
          <h1>
            {isUserEmpty(user) ? "You are NOT signed in" : "You are signed In"}
          </h1>
        </div>{" "}
      </main>
    </>
  );
};

Home.getInitialProps = async ({ req }) => {
  const data = await buildClient({
    host: req?.headers.host,
    cookie: req?.headers.cookie,
  })
    .get("api/users/currentuser")
    .catch((error: HTTPError) => {
      if (error.response?.status === 401) {
        return;
      }
      throw error;
    });
  return data ? UserSchema.parse(await data.json()) : {};
};

export default Home;
