const emptyFunc = () => {}

export default function(func) {
  return func || emptyFunc
}
