import { z } from "zod";

export const UserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  iat: z.number(),
});

export type User = z.infer<typeof UserSchema>;

export const isUserEmpty = (user: User | {}) => {
  return Object.keys(user).length === 0;
};
