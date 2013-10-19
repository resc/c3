c3
==

The common container collection a.k.a. c3 for go.


Introduction
============

This library provides a few basic container interfaces that are missing in Go
in my humble opinion. The slice and map containers are nice, but they don't provide
very many convenience methods. This library aims to remedy that.

Code Quality
============

This library has started its life just a few days ago.
Its not completed yet, and has not seen much of the world either.

Not all containers have implementations, and the query api is not complete,
but I don't expect many changes in the interfaces.

Contributions And Bug Reports
=============================

Send pull requests! 

If you have a bug to report I would be very grateful If
you could submit it as a pull request with a failing test.


Road Map
========

The interfaces Bag, List, Set, Queue and Stack are pretty much complete. 
I don't plan on adding more container interfaces.

There is work to be done in the container implementations, query api, 
and in conversion from and to slices, maps and channels.

Containers
==========

 - Iterable: A container that provides a way to iterate over its elements
 - Bag: An unordered container.
 - List: An ordered, indexable container.
 - Set: An unordered container that does not allow duplicates, and provides a few basic set operations.
 - Stack: A Last-In First-Out container.
 - Queue: A First-In First-Out container.

Thread And Goroutine Safety
===========================

These containers are _not_ goroutine- or thread safe.

There is a basic check for concurrent modification while iterating over a container
so you'll at least know when things get corrupted because of concurrent modifications
Every mutation of a container increments the container version, and Iterator checks
this version on every MoveNext(), and panics if it is not what it expects.

This also means you cannot modify a container while iterating over it.

Example:

	// Don't do this!
	l := c3.ListOf(1,2,3,4)
	for i := l.Iterator(); i.MoveNext(); {
		// vv bug here, cannot modify container while iterating over it! vv
		l.Add(i.Value().(int)*2)
	}

Element Types
=============

Because Go does not have generics (yet..) I choose <code>interface{}</code> for the element type.
This means there will be casting involved with the use of the containers in c3
but Go's type assertions are nice enough to make it only a minor annoyance.

Quering containers
==================

The c3.Query() function provides an entrypoint for the query api of c3.
This api is modelled after the C# Linq api.

Example:

	l := c3.ListOf(1, 2, 3, 4)
	q := c3.Query(l)
	result := q.Where(func(e interface{})) { return e.(int)%2 == 0 }).
	            Select(func(e interface{}) interface{} { return e.(int) * 10 }).
				ToList()

As you can see this api would be much nicer if Go had lambda expressions so that you could type

    e => e.(int) * 10
	
instead of

	func(e interface{}) interface{} { return e.(int) * 10 }
	
as an alternative you can define named functions to improve readability

Example:

	func isMod2(e interface{}) bool {
		return e.(int)%2 == 0
	} 
	func times10(e interface{}) interface{} {
		return e.(int)*10
	} 

	l := c3.ListOf(1, 2, 3, 4)
	q := c3.Query(l)
	result := q.Where(isMod2).
	            Select(times10).
				ToList()	










