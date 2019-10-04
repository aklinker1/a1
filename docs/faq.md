<img width="200" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png"> &emsp;__FAQ__

<details><summary><h2>1. Is this an ORM?</h2></summary>
__No, A1 does not do any database interaction__. All database interactions go though a `DataLoader`, while A1 simply tells the driver what it wants done. This also means that A1 does not handle database setup or teardown. You qwill have to create the tables and manage migrations.
</details>

## 2. Can I still customize a `selectOne` query or any other queries where I don't want the default behavior?

Of course! Check out [this page]() to find out how to override any default behaviors/

## 3. Do you support subscriptions?

As of now, no.
