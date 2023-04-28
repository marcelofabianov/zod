import * as z from 'zod';

const PersonSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  age: z.number().int().positive()
});

type Person = z.infer<typeof PersonSchema>;

const person: Person = PersonSchema.parse({
  "name": "John Doe",
  "email": "john",
  "age": 42
})

console.log(person)
