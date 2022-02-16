import inspect

def foo(a, b, x='blah'):
    pass

print(len(inspect.signature(foo).parameters))

a= [1,2,3]
print(a[:1])
# ArgSpec(args=['a', 'b', 'x'], varargs=None, keywords=None, defaults=('blah',))