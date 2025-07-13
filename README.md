# pokedex

Interactive cli for exploring & catching pokemons through [PokeAPI](https://pokeapi.co)

# Usage
 
## `help`
```bash
    Pokedex > help
    Welcome to the Pokedex!
    Usage:
       map - Print pokemon location areas
       mapb - Print previous pokemon location areas
       explore - Explore pokemons in the location are
       catch - Catch pokemons
       inspect - Inspect pokemons
       pokedex - list of all the names of the Pokemons that were caught
       help - Displays a help message
       exit - Exits from program
```
## `map`
```bash
    Pokedex > map
    mt-coronet-1f-route-216
    mt-coronet-1f-route-211
    mt-coronet-b1f
    great-marsh-area-1
    great-marsh-area-2
    great-marsh-area-3
    great-marsh-area-4
    great-marsh-area-5
    great-marsh-area-6
    solaceon-ruins-2f
    solaceon-ruins-1f
    solaceon-ruins-b1f-a
    solaceon-ruins-b1f-b
    solaceon-ruins-b1f-c
    solaceon-ruins-b2f-a
    solaceon-ruins-b2f-b
    solaceon-ruins-b2f-c
    solaceon-ruins-b3f-a
    solaceon-ruins-b3f-b
```

## `mapb`
```bash
    Pokedex > mapb
    canalave-city-area
    eterna-city-area
    pastoria-city-area
    sunyshore-city-area
    sinnoh-pokemon-league-area
    oreburgh-mine-1f
    oreburgh-mine-b1f
    valley-windworks-area
    eterna-forest-area
    fuego-ironworks-area
    mt-coronet-1f-route-207
    mt-coronet-2f
    mt-coronet-3f
    mt-coronet-exterior-snowfall
    mt-coronet-exterior-blizzard
    mt-coronet-4f
    mt-coronet-4f-small-room
    mt-coronet-5f
    mt-coronet-6f
    mt-coronet-1f-from-exterior
```
## `explore`
```bash
    Pokedex > explore eterna-forest-area
    Exploring eterna-forest-area...
    Found pokemon:
     - caterpie
     - metapod
     - weedle
     - kakuna
     - gastly
     - hoothoot
     - murkrow
     - misdreavus
     - pineco
     - wurmple
     - silcoon
     - beautifly
     - cascoon
     - dustox
     - seedot
     - slakoth
     - nincada
     - bidoof
     - kricketot
     - budew
     - buneary 
 ```

 ## `catch`
 ```bash
    Pokedex > catch pikachu
    Throwing a Pokeball at pikachu...
    pikachu escaped!
    Pokedex > catch pikachu
    Throwing a Pokeball at pikachu...
    pikachu was caught!
    You may now inspect it with the inspect command.
    Pokedex > catch slowpoke 
    Throwing a Pokeball at slowpoke...
    slowpoke escaped!
    Pokedex > catch slowpoke
    Throwing a Pokeball at slowpoke...
    slowpoke was caught!
    You may now inspect it with the inspect command.
    Pokedex > catch squirtle
    Throwing a Pokeball at squirtle...
    squirtle was caught!
    You may now inspect it with the inspect command.
```

## `inspect`
```bash
    Pokedex > inspect pikachu
    Name: pikachu
    Height: 4
    Weight: 60
    Stats:
      -hp: 35
      -attack: 55
      -defense: 40
      -special-attack: 50
      -special-defense: 50
      -speed: 90
    Types:
      - electric
```

## `pokedex`
```bash
    Pokedex > pokedex
    Your Pokedex:
      - squirtle
      - pikachu
      - slowpoke
```