for file in steps/*.dot; do
    dot -Tpng "$file" -o "${file%.dot}.png"
done
