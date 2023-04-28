import * as z from 'zod';

const PersonSchema = z.object({
  name: z.string().trim().min(5).max(255),
  email: z.string().trim().email().min(5).max(255),
  age: z.number().int().positive().min(18).max(120),
});

type Person = z.infer<typeof PersonSchema>;

const person: Person = PersonSchema.parse({
  "name": "John Doe",
  "email": "john",
  "age": 42
})

console.log(person)
