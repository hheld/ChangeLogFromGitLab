[Inspiration for this library](https://github.com/conventional-changelog/conventional-changelog/blob/a5505865ff3dd710cf757f50530e73ef0ca641da/conventions/angular.md)

To enforce the message format, you can use this script on the server as a custom update hook which is a modified version of what you find [here](https://git-scm.com/book/be/v2/Customizing-Git-An-Example-Git-Enforced-Policy):

```ruby
#!/usr/bin/env ruby

$refname = ARGV[0]
$oldrev  = ARGV[1]
$newrev  = ARGV[2]

puts "Enforcing Policies..."
puts "(#{$refname}) (#{$oldrev[0,6]}) (#{$newrev[0,6]})"

$regex = /(?:(feat|fix|docs|style|refactor|perf|test|chore)\(([^\(\)]+)\): ([^\n]+)$\n^$\n((?:\n|.)+)^$\n((?:(?:[Rr]efs|[Cc]loses) #\d+\n)+)|Merge .*)$/m

# enforced custom commit message format
def check_message_format
  missed_revs = `git rev-list #{$oldrev}..#{$newrev}`.split("\n")
  missed_revs.each do |rev|
    message = `git cat-file commit #{rev} | sed '1,/^$/d'`
    if !$regex.match(message)
      puts "[POLICY] Your message is not formatted correctly"
      exit 1
    end
  end
end
check_message_format
```
