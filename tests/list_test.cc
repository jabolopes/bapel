#include "tests/test_util.h"
#include "ts/list.h"

#include <vector>

TEST(ListTest, BasicOperations) {
  ts::List<int> l;
  EXPECT_TRUE(l.empty());
  EXPECT_EQ(l.size(), 0);

  auto l1 = l.add(1);
  EXPECT_FALSE(l1.empty());
  EXPECT_EQ(l1.size(), 1);
  EXPECT_EQ(l1.front(), 1);

  auto l2 = l1.add(2);
  EXPECT_FALSE(l2.empty());
  EXPECT_EQ(l2.size(), 2);
  EXPECT_EQ(l2.front(), 2);

  auto l3 = l2.add(3);
  EXPECT_FALSE(l3.empty());
  EXPECT_EQ(l3.size(), 3);
  EXPECT_EQ(l3.front(), 3);

  auto l4 = l3.remove();
  EXPECT_EQ(l4.size(), 2);
  EXPECT_EQ(l4.front(), 2);

  auto l5 = l4.remove();
  EXPECT_EQ(l5.size(), 1);
  EXPECT_EQ(l5.front(), 1);

  auto l6 = l5.remove();
  EXPECT_TRUE(l6.empty());
  EXPECT_EQ(l6.size(), 0);

  auto l7 = l6.remove();
  EXPECT_TRUE(l7.empty());
  EXPECT_EQ(l7.size(), 0);
}

TEST(ListTest, IterationAndCollect) {
  ts::List<int> l;
  l = l.add(1).add(2).add(3);

  // l.to_vector() collects forward
  std::vector<int> want = {1, 2, 3};
  EXPECT_EQ(l.to_vector(), want);

  // iterate via next() / linked list
  std::vector<int> reverse_order;
  for (auto curr = l; !curr.empty(); curr = curr.remove()) {
    reverse_order.push_back(curr.front());
  }
  std::vector<int> want_reverse = {3, 2, 1};
  EXPECT_EQ(reverse_order, want_reverse);
}

TEST(ListTest, Reverse) {
  ts::List<int> l;
  l = l.add(10).add(20).add(30);

  auto r = l.reverse();
  std::vector<int> want_r = {30, 20, 10};
  EXPECT_EQ(r.to_vector(), want_r);
}
