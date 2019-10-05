# <img height="25" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png"> &ensp;FAQ

### 1. Is this an ORM?

<details>
<summary>No...</summary>

__`a1` does not do any database interaction__. All database interactions go though a `DataLoader`, while `a1` simply tells the data loader what it wants done. This also means that `a1` does not handle database setup or teardown. You will have to create the tables and manage migrations.
    
</details>

<br/>

### 2. Can I still customize a `GetOne` query or any other queries where I don't want the default behavior?

<details>
<summary>Of course!</summary>

Of course! Check out [this page]() to find out how to override any default behaviors.

</details>

<br/>

### 3. Do you support subscriptions?

<details>
<summary>Nope.</summary>

There is no support for it currently, and no interest in doing so in the future.

</details>
